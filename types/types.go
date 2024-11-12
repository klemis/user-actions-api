package types

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Action struct {
	ID         int       `json:"id"`
	Type       string    `json:"type"` // use type
	UserID     int       `json:"userId"`
	TargetUser int       `json:"targetUser"`
	CreatedAt  time.Time `json:"createdAt"`
}

// ActionsProbalibity holds the probability for each possible next action.
type ActionsProbalibity map[string]float64

// Referral represents mapping of users to the IDs of users they referred.
type Referral map[int][]int

// ReferralIndex store the referral index for each user.
type ReferralIndex map[int]int
