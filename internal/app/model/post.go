package model

import "time"

type Post struct {
	Score            int        `json:"score"`
	Views            int        `json:"views"`
	Type             string     `json:"type"`
	Title            string     `json:"title"`
	Author           *User      `json:"author"`
	Category         string     `json:"category"`
	Text             string     `json:"text"`
	Votes            []*Vote    `json:"votes"`
	Comments         []*Comment `json:"comments"`
	Created          time.Time  `json:"created"`
	UpvotePercentage int        `json:"upvotePercentage"`
	ID               string     `json:"id"`
}

// {
//     "score": 0,
//     "views": 28,
//     "type": "text",
//     "title": "Hello",
//     "author": {
//         "username": "diazabdulm",
//         "id": "5da016d78036b7fc8bdf0526"
//     },
//     "category": "videos",
//     "text": "World",
//     "votes": [{
//         "user": "5da016d78036b7fc8bdf0526",
//         "vote": 1
//     }, {
//         "user": "5e3bc44b0baa2100072a7051",
//         "vote": -1
//     }],
//     "comments": [{
//         "created": "2019-10-18T12:33:02.758Z",
//         "author": {
//             "username": "hhfyhufy7895t77",
//             "id": "5da9b05ae09fe7cb6259a213"
//         },
//         "body": "ig kjgjhf fhifguiguiugiuguig",
//         "id": "5da9b0fee09fe7bdad59a216"
//     }],
//     "created": "2019-10-11T05:45:13.020Z",
//     "upvotePercentage": 50,
//     "id": "5da016e98036b76ef6df0527"
// }
