package main

// import "cloud"
import "math"


// VPC component
#VPC: {
  name:     string
  cidr:     string
}

// Subnet component
#Subnet: {
  name:     string
  cidr:     string
  // zone:     string
  vpc:      #VPC
}

// Storage bucket
#Bucket: {
  name:     string
  location: string
  labels: [string]: string
}


let myVPC = #VPC & {
  name: "main-vpc"
  cidr: "10.0.0.0/16"
}

let subnet = #Subnet & {
  name: "public-subnet-a"
  cidr: "10.0.1.0/24"
  vpc: myVPC
}

let bucket = #Bucket & {
  name: "app-logs"
  location: "lol"
  labels: {
    env: "prod"
    network: myVPC.name
  }
}

foo: "bar"

config: {
  resources: {
    Vpc: myVPC,
    Subnet: subnet,
    Bucket: bucket
  },
  foo2: foo
}
bar: foo
bar2: "yeee" + (foo + "asd" + bar)

_string1: "foo"
_string2: "bar"
_int1: 1
_int2: 2
_float1: 1.0
_float2: 2.0
_bool1: true
_bool2: false
_list: [1,2,3,4,5,6,7,8,9,10]
_struct1: {
  string1: "foo"
  int1: _int1
  float1: 1.0
  bool1: true
  list1: [1,2,_int2,4]
}
_struct2: {
  string2: "foo"
  int2: _int2
  float2: 2.0
  bool2: true
  list2: [2,2,_int2,4]
}

// NoOp:
AndOp: _bool1 & _bool1
OrOp: _bool1 | _bool2
SelectorOp: _string1
IndexOp: _list[_int1]
SliceOp: _list[_int1:_int2]
CallOp: math.Atan2(_int1, _int2)
BooleanAndOp: _bool1 && _bool2
BooleanOrOp: _bool1 || _bool2
EqualOp: _int1 == _int2
NotOp: !_bool1
NotEqualOp: _int1 != _int2
LessThanOp: _int1 < _int2
LessThanEqualOp: _int1 <= _int2
GreaterThanOp: _int1 > _int2
GreaterThanEqualOp: _int1 >= _int2
RegexMatchOp: _string1 =~ _string2
NotRegexMatchOp: _string1 !~ _string2
AddOp: _int1 + _int2
SubtractOp: _int1 - _int2
MultiplyOp: _int1 * _int2
FloatQuotientOp: _int1 / _int2
// IntQuotientOp: // Ich glaub das wird aktuell nicht genutzt.
// IntRemainderOp: // Ich glaub das wird aktuell nicht genutzt.
// IntDivideOp: // Ich glaub das wird aktuell nicht genutzt.
// IntModuloOp: // Ich glaub das wird aktuell nicht genutzt.
InterpolationOp: "Hello \(_string1)!"
