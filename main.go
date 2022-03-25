package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//DB store the database session information. Needs to be initialized once type DB struct

type DB struct {
	session *mgo.Session
	collection *mgo.Collection
}



// Movie hold a movie data

 type Movie struct {
	 ID  	bson.ObjectId `json:"id" bson:"_id, omitempty"`
	 Name	string `json:"name" bson:"name"`
	 Year 	string `json:"year" bson:"year"`
	 Directors []string 	`json:"directors" bson:"directors"`
	 Writes []string 	`json:"writers" bson:"writes"`
	 BoxOffice  BoxOffice	`json:"boxOffice" bson:"boxOffice"`
 }

 //BoxOffice is nested in Movie

 type BoxOffice struct {
	 Budget uint64 `json:"budget" bson:"budget"`
	 Gross uint64 `json:"gross" bson:"gross"`
 }
 
 //GetMovie fetches a movie with a given ID

 func (db *DB) GetMovie(w http.ResponseWriter, r *http.Request){
	 vars := mux.Vars(r)

	 w.WriteHeader(http.StatusOK)
	 var movie Movie
	 err := db.collection.Find(bson.M{"_id" : bson.ObjectIdHex(vars["id"])}).One(&movie)
	 if err != nil {
		 w.Write([]byte(err.Error()))
	 } else {
		 w.Header().Set("Content-Type", "application/json")
		 response, _ := json.Marshal(movie)
		 w.Write(response)
	 }
 }

//PostMovie adds a new movie to  our MongoDB collection

func (db *DB) PostMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	postBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(postBody, &movie)
	// Create a hash ID to insert
	movie.ID = bson.NewObjectId()
	err := db.collection.Insert(movie)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _:= json.Marshal(movie)
		w.Write(response)
	}
}

//UpdateMovie modifies the data of given resource

func (db *DB) UpdateMovie(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	var movie Movie
	putBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(putBody, &movie)
	//Create an Hash ID to insert
	err := db.collection.Update(bson.M{"_id": bson.IsObjectIdHex(vars["id"])}, bson.M{"$set": &movie})
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text")
		w.Write([]byte("Update succesfully!"))
	}
}

//DeleteMovie removes the data from the db
func (db *DB) DeleteMovie(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	//Create an Hash ID to insert
	err := db.collection.Remove(bson.M{"_id":bson.IsObjectIdHex(vars["id"])})
	if err != nil  {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text")
		w.Write([]byte("Delete succesfully"))
	}
}

func main(){
	session, err := mgo.Dial("127.0.0.1")
	c := session.DB("appdb").C("movies")
	db := &DB{session: session, collection: c}
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//Create a new router
	r := mux.NewRouter()
	//Attach an elegant path with handler
	r.HandleFunc("/v1/movies/{id:[a-zA-Z0-9]*}", db.GetMovie).Methods("GET")
	r.HandleFunc("/v1/movies", db.PostMovie).Methods("POST")
	srv := &http.Server{
		Handler: r,
		Addr: "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
