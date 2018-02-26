package main

import (
	"os"
	"net/http"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	)

type KeyValue struct {
	Key string
	Value string
}

func setupRouter() *gin.Engine {
	var m map[string]string
	m = make(map[string]string)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/v1/keyvalue/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, ok := m[key]
		if (ok) {
			c.String(http.StatusOK, value)
		} else {
			c.String(http.StatusNotFound, "not found")
		}
	})
	r.POST("/v1/keyvalue", func(c *gin.Context) {
		key := c.Query("key")
		value := c.Query("value")
		m[key] = value
		c.String(http.StatusOK, "ok: value set")
	})

	return r
}

func main() {
	r := setupRouter()
	mongoUrl := os.Args[1]
	log.Println("mongo URL: " + mongoUrl)
	session, err := mgo.Dial(mongoUrl)
	if err == nil {
		mongoDatabase := os.Args[2]
		mongoCollectionName := os.Args[3]
		log.Println("mongo Database: " + mongoDatabase)
		log.Println("mongo Collection: " + mongoCollectionName)
		mongoCollection := session.DB(mongoDatabase).C(mongoCollectionName)
		index := mgo.Index{
	        Key:        []string{"key"},
	        Unique:     true,
	        DropDups:   true,
	        Background: true,
	        Sparse:     true,
	    }
	    err = mongoCollection.EnsureIndex(index)
	    if err != nil {
	        log.Println(err)
	    }
		r.GET("/v2/keyvalue/:key", func(c *gin.Context) {
			key := c.Param("key")
			result := KeyValue{}
			mongoCollection.Find(bson.M{"key": key}).One(&result)
			if result.Key == key {
				c.String(http.StatusOK, result.Value)
			} else {
				c.String(http.StatusNotFound, "not found")
			}
		})
		r.POST("/v2/keyvalue", func(c *gin.Context) {
			key := c.Query("key")
			value := c.Query("value")
			keyValue := KeyValue{key, value}
			result := KeyValue{}
			mongoCollection.Find(bson.M{"key": key}).One(&result)
			if result.Key == key {
				mongoCollection.Update(bson.M{"key": key}, &keyValue)
			} else {
				mongoCollection.Insert(keyValue)
			}
			c.String(http.StatusOK, "ok: value set")
		})
	} else {
		log.Fatal("Can't start mongodb version")
		log.Fatal(err)
	}

	r.Run(":8888") // listen and serve on 0.0.0.0:8888
}