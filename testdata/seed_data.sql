-- Seed data for movies database
-- This file is automatically executed when PostgreSQL starts in Docker

-- Wait for the movies table to be created by migrations
-- Note: In production, this would be handled differently

-- Insert sample movies
INSERT INTO movies (title, director, year, genre, rating, description, duration, language, country) VALUES
('The Shawshank Redemption', 'Frank Darabont', 1994, ARRAY['Drama'], 9.3, 'Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.', 142, 'English', 'USA'),
('The Godfather', 'Francis Ford Coppola', 1972, ARRAY['Crime', 'Drama'], 9.2, 'The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.', 175, 'English', 'USA'),
('The Dark Knight', 'Christopher Nolan', 2008, ARRAY['Action', 'Crime', 'Drama'], 9.0, 'When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.', 152, 'English', 'USA'),
('Pulp Fiction', 'Quentin Tarantino', 1994, ARRAY['Crime', 'Drama'], 8.9, 'The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.', 154, 'English', 'USA'),
('The Lord of the Rings: The Return of the King', 'Peter Jackson', 2003, ARRAY['Adventure', 'Drama', 'Fantasy'], 8.9, 'Gandalf and Aragorn lead the World of Men against Sauron''s army to draw his gaze from Frodo and Sam as they approach Mount Doom with the One Ring.', 201, 'English', 'New Zealand'),
('Forrest Gump', 'Robert Zemeckis', 1994, ARRAY['Drama', 'Romance'], 8.8, 'The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man with an IQ of 75.', 142, 'English', 'USA'),
('Inception', 'Christopher Nolan', 2010, ARRAY['Action', 'Sci-Fi', 'Thriller'], 8.8, 'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.', 148, 'English', 'USA'),
('The Matrix', 'Lana Wachowski, Lilly Wachowski', 1999, ARRAY['Action', 'Sci-Fi'], 8.7, 'A computer programmer discovers that reality as he knows it is a simulation created by machines to distract humans while using their bodies as an energy source.', 136, 'English', 'USA'),
('Goodfellas', 'Martin Scorsese', 1990, ARRAY['Biography', 'Crime', 'Drama'], 8.7, 'The story of Henry Hill and his life in the mob, covering his relationship with his wife Karen Hill and his mob partners Jimmy Conway and Tommy DeVito.', 146, 'English', 'USA'),
('Star Wars: Episode V - The Empire Strikes Back', 'Irvin Kershner', 1980, ARRAY['Action', 'Adventure', 'Fantasy'], 8.7, 'After the Rebels are brutally overpowered by the Empire on the ice planet Hoth, Luke Skywalker begins Jedi training with Yoda.', 124, 'English', 'USA'),
('Parasite', 'Bong Joon Ho', 2019, ARRAY['Comedy', 'Drama', 'Thriller'], 8.6, 'Greed and class discrimination threaten the newly formed symbiotic relationship between the wealthy Park family and the destitute Kim clan.', 132, 'Korean', 'South Korea'),
('Spirited Away', 'Hayao Miyazaki', 2001, ARRAY['Animation', 'Adventure', 'Family'], 8.6, 'During her family''s move to the suburbs, a sullen 10-year-old girl wanders into a world ruled by gods, witches, and spirits.', 125, 'Japanese', 'Japan'),
('The Silence of the Lambs', 'Jonathan Demme', 1991, ARRAY['Crime', 'Drama', 'Thriller'], 8.6, 'A young F.B.I. cadet must receive the help of an incarcerated and manipulative cannibal killer to help catch another serial killer.', 118, 'English', 'USA'),
('Saving Private Ryan', 'Steven Spielberg', 1998, ARRAY['Drama', 'War'], 8.6, 'Following the Normandy Landings, a group of U.S. soldiers go behind enemy lines to retrieve a paratrooper whose brothers have been killed in action.', 169, 'English', 'USA'),
('Schindler''s List', 'Steven Spielberg', 1993, ARRAY['Biography', 'Drama', 'History'], 8.9, 'In German-occupied Poland during World War II, industrialist Oskar Schindler gradually becomes concerned for his Jewish workforce after witnessing their persecution by the Nazis.', 195, 'English', 'USA'),
('12 Angry Men', 'Sidney Lumet', 1957, ARRAY['Crime', 'Drama'], 9.0, 'A jury holdout attempts to prevent a miscarriage of justice by forcing his colleagues to reconsider the evidence.', 96, 'English', 'USA'),
('Fight Club', 'David Fincher', 1999, ARRAY['Drama'], 8.8, 'An insomniac office worker and a devil-may-care soapmaker form an underground fight club that evolves into something much, much more.', 139, 'English', 'USA'),
('The Good, the Bad and the Ugly', 'Sergio Leone', 1966, ARRAY['Western'], 8.8, 'A bounty hunting scam joins two men in an uneasy alliance against a third in a race to find a fortune in gold buried in a remote cemetery.', 178, 'Italian', 'Italy'),
('City of God', 'Fernando Meirelles', 2002, ARRAY['Crime', 'Drama'], 8.6, 'In the slums of Rio, two kids'' paths diverge as one struggles to become a photographer and the other a kingpin.', 130, 'Portuguese', 'Brazil'),
('Se7en', 'David Fincher', 1995, ARRAY['Crime', 'Drama', 'Mystery'], 8.6, 'Two detectives, a rookie and a veteran, hunt a serial killer who uses the seven deadly sins as his motives.', 127, 'English', 'USA');

-- Note: Poster data would be added later through the application
-- For testing purposes, we're not including actual image data in the seed