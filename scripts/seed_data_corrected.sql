-- Movies MCP Server - Sample Seed Data
-- This script populates the database with sample movie data for testing and development

-- Clear existing data (optional)
-- TRUNCATE movies RESTART IDENTITY CASCADE;

-- Insert sample movies data
INSERT INTO movies (title, director, year, genre, rating, description) VALUES
    ('The Shawshank Redemption', 'Frank Darabont', 1994, '{"Drama"}', 9.3, 'Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.'),
    
    ('The Godfather', 'Francis Ford Coppola', 1972, '{"Crime","Drama"}', 9.2, 'The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.'),
    
    ('The Dark Knight', 'Christopher Nolan', 2008, '{"Action","Crime","Drama"}', 9.0, 'When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests.'),
    
    ('The Godfather Part II', 'Francis Ford Coppola', 1974, '{"Crime","Drama"}', 9.0, 'The early life and career of Vito Corleone in 1920s New York City is portrayed, while his son, Michael, expands and tightens his grip on the family crime syndicate.'),
    
    ('12 Angry Men', 'Sidney Lumet', 1957, '{"Drama"}', 9.0, 'A jury holdout attempts to prevent a miscarriage of justice by forcing his colleagues to reconsider the evidence.'),
    
    ('Schindlers List', 'Steven Spielberg', 1993, '{"Biography","Drama","History"}', 8.9, 'In German-occupied Poland during World War II, industrialist Oskar Schindler gradually becomes concerned for his Jewish workforce after witnessing their persecution.'),
    
    ('The Lord of the Rings: The Return of the King', 'Peter Jackson', 2003, '{"Adventure","Drama","Fantasy"}', 8.9, 'Gandalf and Aragorn lead the World of Men against Saurons army to draw his gaze from Frodo and Sam as they approach Mount Doom with the One Ring.'),
    
    ('Pulp Fiction', 'Quentin Tarantino', 1994, '{"Crime","Drama"}', 8.9, 'The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.'),
    
    ('The Good, the Bad and the Ugly', 'Sergio Leone', 1966, '{"Western"}', 8.8, 'A bounty hunting scam joins two men in an uneasy alliance against a third in a race to find a fortune in gold buried in a remote cemetery.'),
    
    ('Fight Club', 'David Fincher', 1999, '{"Drama"}', 8.8, 'An insomniac office worker and a devil-may-care soap maker form an underground fight club that evolves into an anarchist organization.'),
    
    ('Forrest Gump', 'Robert Zemeckis', 1994, '{"Drama","Romance"}', 8.8, 'The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man.'),
    
    ('Inception', 'Christopher Nolan', 2010, '{"Action","Sci-Fi","Thriller"}', 8.7, 'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.'),
    
    ('The Lord of the Rings: The Two Towers', 'Peter Jackson', 2002, '{"Adventure","Drama","Fantasy"}', 8.7, 'While Frodo and Sam edge closer to Mordor with the help of the shifty Gollum, the divided fellowship makes a stand against Saurons new ally, Saruman.'),
    
    ('Star Wars: Episode V - The Empire Strikes Back', 'Irvin Kershner', 1980, '{"Action","Adventure","Fantasy","Sci-Fi"}', 8.7, 'After the Rebels are brutally overpowered by the Empire on the ice planet Hoth, Luke Skywalker begins Jedi training with Yoda.'),
    
    ('The Lord of the Rings: The Fellowship of the Ring', 'Peter Jackson', 2001, '{"Adventure","Drama","Fantasy"}', 8.6, 'A meek Hobbit from the Shire and eight companions set out on a journey to destroy the powerful One Ring and save Middle-earth from the Dark Lord Sauron.'),
    
    ('Goodfellas', 'Martin Scorsese', 1990, '{"Biography","Crime","Drama"}', 8.6, 'The story of Henry Hill and his life in the mob, covering his relationship with his wife Karen Hill and his mob partners Jimmy Conway and Tommy DeVito.'),
    
    ('One Flew Over the Cuckoos Nest', 'Milos Forman', 1975, '{"Drama"}', 8.6, 'A criminal pleads insanity and is admitted to a mental institution, where he rebels against the oppressive nurse and rallies up the scared patients.'),
    
    ('The Matrix', 'The Wachowskis', 1999, '{"Action","Sci-Fi"}', 8.6, 'When a beautiful stranger leads computer hacker Neo to a forbidding underworld, he discovers the shocking truth--the life he knows is the elaborate deception of an evil cyber-intelligence.'),
    
    ('Seven Samurai', 'Akira Kurosawa', 1954, '{"Action","Adventure","Drama"}', 8.6, 'A poor village under attack by bandits recruits seven unemployed samurai to help them defend themselves.'),
    
    ('City of God', 'Fernando Meirelles', 2002, '{"Crime","Drama"}', 8.6, 'In the slums of Rio, two kids paths diverge as one struggles to become a photographer and the other a kingpin.');

-- Display statistics after insertion
SELECT 
    COUNT(*) as total_movies,
    MIN(year) as earliest_year,
    MAX(year) as latest_year,
    ROUND(AVG(rating), 2) as average_rating,
    COUNT(DISTINCT unnest(genre)) as unique_genres,
    COUNT(DISTINCT director) as unique_directors
FROM movies;