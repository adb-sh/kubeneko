package main

// import "cloud"


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


config: {
  resources: {
    Vpc: myVPC,
    Subnet: subnet,
    Bucket: bucket
  },
  foo2: nested.foo
}

nested: {
  foo: "bar"
}

bar: nested.foo
bar2: "yeee" + (nested.foo + "asd") + bar
