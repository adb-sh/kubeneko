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

foo: baa: "bar"

config: {
  resources: [
    myVPC,
    subnet,
    bucket
  ],
  foo2: foo
}
bar: foo.baa
bar2: foo.baa + "asd"


foo: asd: "dsa"
foo2: foo.asd
foo3: foo2 + "-sfdj" + foo.asd
