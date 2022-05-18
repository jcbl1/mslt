package main

import (
	"os"
	"regexp"
	"strings"
	"time"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Work struct {
	Artist, Name string
}
type MusicList struct {
	Date string
	List []Work
}

func main() {
	f, err := os.Open("/home/jeff/Documents/tmp/musiclist")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		panic(err.Error())
	}

	b := make([]byte, stat.Size())

	_, err = f.Read(b)
	if err != nil {
		panic(err.Error())
	}

	pileTr := regexp.MustCompile(`<tr.*?</tr>`)
	pileTt := regexp.MustCompile(`<b title.*?</div`)
	pileArtist := regexp.MustCompile(`<span title.*?>`)
	pileQuo := regexp.MustCompile(`".*?"`)
	res := pileTr.FindAll(b, -1)

	musicList := MusicList{Date: time.Now().Format("2006-1-2")}

	for _, resn := range res {
		resTt := pileTt.Find(resn)
		resArtist := pileArtist.Find(resn)
		musicList.List = append(musicList.List, Work{
			strings.Trim(string(pileQuo.Find(resArtist)), "\""),
			strings.Trim(strings.ReplaceAll(strings.ReplaceAll(string(pileQuo.Find(resTt)), "&nbsp;", " "), "&quot;", `"`), "\""),
		})
	}

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://jeff:1234@localhost:27017/jeffDb"))
	if err != nil {
		panic(err.Error())
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		panic(err.Error())
	}
	defer mongoClient.Disconnect(ctx)

	collection := mongoClient.Database("jeffDb").Collection("musicDaily")

	_, err = collection.InsertOne(context.TODO(), musicList)
	if err != nil {
		panic(err.Error())
	}
}
