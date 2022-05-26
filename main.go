package main

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlayerInfo struct {
	Id       int
	NickName string

	CurrentRoomId int

	CurrentCard string
	LastCard    string

	SkipCount int
}

type RoomInfo struct {
	Id int

	MasterPlayer int

	PlayerInfos []PlayerInfo
}

var RoomInfos map[int]RoomInfo

var PlayerInfos map[int]PlayerInfo

func main() {
	r := gin.Default()

	RoomInfos = make(map[int]RoomInfo, 0)
	PlayerInfos = make(map[int]PlayerInfo, 0)

	type LoginMessage struct {
		NickName string
	}
	r.GET("/login", func(c *gin.Context) {
		req := LoginMessage{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "Bad Request", "error": err.Error()})
			return
		}

		playerInfo := PlayerInfo{
			Id:       rand.Int(),
			NickName: req.NickName,
		}

		PlayerInfos[playerInfo.Id] = playerInfo

		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "playerId": playerInfo.Id})
	})

	type PlayerIdMessage struct {
		PlayerId int
	}
	r.GET("/createroom", func(c *gin.Context) {
		req := PlayerIdMessage{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "Bad Request", "error": err.Error()})
			return
		}

		playerInfo := PlayerInfos[req.PlayerId]
		roomId := rand.Int()
		playerInfo.CurrentRoomId = roomId

		roomInfo := RoomInfo{
			Id:           roomId,
			MasterPlayer: req.PlayerId,
			PlayerInfos:  []PlayerInfo{playerInfo},
		}

		RoomInfos[roomInfo.Id] = roomInfo

		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "roomId": roomInfo.Id})
	})

	type roomPlayerRequest struct {
		RoomId   int
		PlayerId int
	}
	r.GET("/joinroom", func(c *gin.Context) {
		req := roomPlayerRequest{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "Bad Request", "error": err.Error()})
			return
		}

		roomInfo := RoomInfos[req.RoomId]
		playerInfo := PlayerInfos[req.PlayerId]
		playerInfo.CurrentRoomId = req.RoomId

		roomInfo.PlayerInfos = append(roomInfo.PlayerInfos, playerInfo)

		RoomInfos[req.RoomId] = roomInfo
		PlayerInfos[req.PlayerId] = playerInfo

		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "roomId": req.RoomId})
	})

	r.GET("/leaveroom", func(c *gin.Context) {
		req := roomPlayerRequest{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "Bad Request", "error": err.Error()})
			return
		}

		newPlayerInfos := make([]PlayerInfo, 0)

		roomInfo := RoomInfos[req.RoomId]

		playerInfo := PlayerInfos[req.PlayerId]
		playerInfo.CurrentRoomId = 0
		playerInfo.SkipCount = 0

		for _, p := range roomInfo.PlayerInfos {
			if p.Id != playerInfo.Id {
				newPlayerInfos = append(newPlayerInfos, p)
			}
		}
		roomInfo.PlayerInfos = newPlayerInfos

		if len(newPlayerInfos) != 0 {
			RoomInfos[req.RoomId] = roomInfo
		} else {
			delete(RoomInfos, req.RoomId)
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
	})

	type RoomIdMessage struct {
		RoomId int
	}

	r.GET("/roominfo", func(c *gin.Context) {
		req := RoomIdMessage{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "Bad Request", "error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "roomInfo": RoomInfos[req.RoomId]})
	})

	type RoomInfo struct {
		RoomId     int
		MasterName string
	}
	r.GET("/roomlist", func(c *gin.Context) {
		roomList := make([]RoomInfo, 0)

		for k, v := range RoomInfos {
			roomList = append(roomList, RoomInfo{k, PlayerInfos[v.MasterPlayer].NickName})
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "roomList": roomList})
	})

	r.GET("/getnextcard", func(c *gin.Context) {
		req := PlayerIdMessage{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "Bad Request", "error": err.Error()})
			return
		}

	})

	r.GET("/refresh", func(c *gin.Context) {
		req := roomPlayerRequest{}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": -1, "msg": "Bad Request", "error": err.Error()})
			return
		}

		roomInfo := RoomInfos[req.RoomId]

		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok", "roomInfo": roomInfo})
	})

	panic(r.Run(":8080"))
}
