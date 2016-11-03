use sf-movie-locations;

db.createUser(
  {
    user: "administrator",
    pwd: "superpassword",
    roles: [ { role: "dbAdmin", db: "sf-movie-locations" } ]
  }
);

db.createUser(
  {
    user: "dataloader",
    pwd: "qwerty23",
    roles: [ { role: "readWrite", db: "sf-movie-locations" } ]
  }
);

db.createUser(
  {
    user: "webapp",
    pwd: "asdf88",
    roles: [ { role: "read", db: "sf-movie-locations" } ]
  }
);
