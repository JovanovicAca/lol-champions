CREATE TABLE "champions" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(20) NOT NULL,
  "class" varchar(20) NOT NULL,
  "weapon" varchar(20) NOT NULL,
  "MagicCost" varchar(20) NOT NULL
);

CREATE TABLE "worlds" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(10) NOT NULL,
  "description" varchar(100) NOT NULL
);

CREATE TABLE "position" (
  "id" SERIAL PRIMARY KEY,
  "position" varchar(10) UNIQUE NOT NULL
);

CREATE TABLE "champion_position" (
  "championId" uuid,
  "positionId" uuid,
  PRIMARY KEY ("championId"),
  PRIMARY KEY ("positionId")
);

ALTER TABLE "champion_position" ADD FOREIGN KEY ("championId") REFERENCES "champions" ("id");

ALTER TABLE "champion_position" ADD FOREIGN KEY ("positionId") REFERENCES "position" ("id");

ALTER TABLE "champions" ADD FOREIGN KEY ("id") REFERENCES "worlds" ("id");
