package common

import "errors"

var HashPassword = errors.New("invalid hash format")
var PasswordIncorrect = errors.New("password not correct")
var TokenExpired = errors.New("token is expired")
var TokenInvalid = errors.New("token is invalid")
var TokenMissed = errors.New("token is missing")
var RedisSentinelMaster = errors.New("no valid master found from sentinel")
