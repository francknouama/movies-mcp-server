-- Movies MCP Server - Sample Seed Data
-- This script populates the database with sample movie data for testing and development

-- Clear existing data (optional)
-- TRUNCATE movies RESTART IDENTITY CASCADE;

-- Insert sample movies data
INSERT INTO movies (title, director, year, genre, rating, description) VALUES
    ('The Shawshank Redemption', 'Frank Darabont', 1994, 'Drama', 9.3, 'Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.', 'https://image.tmdb.org/t/p/w500/q6y0Go1tsGEsmtFryDOJo3dEmqu.jpg'),
    
    ('The Godfather', 'Francis Ford Coppola', 1972, 'Crime', 9.2, 'The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.', 'https://image.tmdb.org/t/p/w500/3bhkrj58Vtu7enYsRolD1fZdja1.jpg'),
    
    ('The Dark Knight', 'Christopher Nolan', 2008, 'Action', 9.0, 'When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests.', 'https://image.tmdb.org/t/p/w500/qJ2tW6WMUDux911r6m7haRef0WH.jpg'),
    
    ('The Godfather Part II', 'Francis Ford Coppola', 1974, 'Crime', 9.0, 'The early life and career of Vito Corleone in 1920s New York City is portrayed, while his son, Michael, expands and tightens his grip on the family crime syndicate.', 'https://image.tmdb.org/t/p/w500/hek3koDUyRQk7FIhPXsa6mT2Zc3.jpg'),
    
    ('12 Angry Men', 'Sidney Lumet', 1957, 'Drama', 9.0, 'A jury holdout attempts to prevent a miscarriage of justice by forcing his colleagues to reconsider the evidence.', 'https://image.tmdb.org/t/p/w500/ow3wq89wM8qd5X7hWKxiRfsFf9C.jpg'),
    
    ('Schindler''s List', 'Steven Spielberg', 1993, 'Biography', 8.9, 'In German-occupied Poland during World War II, industrialist Oskar Schindler gradually becomes concerned for his Jewish workforce after witnessing their persecution.', 'https://image.tmdb.org/t/p/w500/sF1U4EUQS8YHUYjNl3pMGNIQyr0.jpg'),
    
    ('The Lord of the Rings: The Return of the King', 'Peter Jackson', 2003, 'Adventure', 8.9, 'Gandalf and Aragorn lead the World of Men against Sauron''s army to draw his gaze from Frodo and Sam as they approach Mount Doom with the One Ring.', 'https://image.tmdb.org/t/p/w500/rCzpDGLbOoPwLjy3OAm5NUPOTrC.jpg'),
    
    ('Pulp Fiction', 'Quentin Tarantino', 1994, 'Crime', 8.9, 'The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.', 'https://image.tmdb.org/t/p/w500/d5iIlFn5s0ImszYzBPb8JPIfbXD.jpg'),
    
    ('The Lord of the Rings: The Fellowship of the Ring', 'Peter Jackson', 2001, 'Adventure', 8.8, 'A meek Hobbit from the Shire and eight companions set out on a journey to destroy the powerful One Ring and save Middle-earth from the Dark Lord Sauron.', 'https://image.tmdb.org/t/p/w500/6oom5QYQ2yQTMJIbnvbkBL9cHo6.jpg'),
    
    ('The Good, the Bad and the Ugly', 'Sergio Leone', 1966, 'Western', 8.8, 'A bounty hunting scam joins two men in an uneasy alliance against a third in a race to find a fortune in gold buried in a remote cemetery.', 'https://image.tmdb.org/t/p/w500/bX2xnavhMYjWDoZp1VM6VnU1xwe.jpg'),
    
    ('Forrest Gump', 'Robert Zemeckis', 1994, 'Drama', 8.8, 'The presidencies of Kennedy and Johnson, the Vietnam War, the Watergate scandal and other historical events unfold from the perspective of an Alabama man.', 'https://image.tmdb.org/t/p/w500/arw2vcBveWOVZr6pxd9XTd1TdQa.jpg'),
    
    ('Fight Club', 'David Fincher', 1999, 'Drama', 8.8, 'An insomniac office worker and a devil-may-care soap maker form an underground fight club that evolves into an anarchist organization.', 'https://image.tmdb.org/t/p/w500/pB8BM7pdSp6B6Ih7QZ4DrQ3PmJK.jpg'),
    
    ('The Lord of the Rings: The Two Towers', 'Peter Jackson', 2002, 'Adventure', 8.7, 'While Frodo and Sam edge closer to Mordor with the help of the shifty Gollum, the divided fellowship makes a stand against Sauron''s new ally, Saruman.', 'https://image.tmdb.org/t/p/w500/5VTN0pR8gcqV3EPUHHfMGnJYN9L.jpg'),
    
    ('Inception', 'Christopher Nolan', 2010, 'Action', 8.7, 'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a C.E.O.', 'https://image.tmdb.org/t/p/w500/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg'),
    
    ('The Empire Strikes Back', 'Irvin Kershner', 1980, 'Adventure', 8.7, 'After the Rebels are brutally overpowered by the Empire on the ice planet Hoth, Luke Skywalker begins Jedi training with Yoda, while his friends are pursued across the galaxy.', 'https://image.tmdb.org/t/p/w500/nNAeTmF4CtdSgMDplXTDPOpYzsX.jpg'),
    
    ('The Matrix', 'The Wachowskis', 1999, 'Action', 8.7, 'A computer programmer discovers that reality as he knows it is not real and finds himself in a war between humanity and machines.', 'https://image.tmdb.org/t/p/w500/f89U3ADr1oiB1s9GkdPOEpXUk5H.jpg'),
    
    ('Goodfellas', 'Martin Scorsese', 1990, 'Biography', 8.7, 'The story of Henry Hill and his life in the mob, covering his relationship with his wife Karen Hill and his mob partners Jimmy Conway and Tommy DeVito.', 'https://image.tmdb.org/t/p/w500/aKuFiU82s5ISJpGZp7YkIr3kCUd.jpg'),
    
    ('One Flew Over the Cuckoo''s Nest', 'Milos Forman', 1975, 'Drama', 8.7, 'A criminal pleads insanity and is admitted to a mental institution, where he rebels against the oppressive nurse and rallies up the scared patients.', 'https://image.tmdb.org/t/p/w500/3jcbDmRFiQ83drXNOvRDeKHxS0V.jpg'),
    
    ('Se7en', 'David Fincher', 1995, 'Crime', 8.6, 'Two detectives, a rookie and a veteran, hunt a serial killer who uses the seven deadly sins as his motives.', 'https://image.tmdb.org/t/p/w500/69Sns8WoET6CfaYlIkHbla4l7nC.jpg'),
    
    ('Seven Samurai', 'Akira Kurosawa', 1954, 'Adventure', 8.6, 'A poor village under attack by bandits recruits seven unemployed samurai to help them defend themselves.', 'https://image.tmdb.org/t/p/w500/8OKmBV5BUFzmozIC3pPWKHy17kx.jpg'),
    
    ('The Silence of the Lambs', 'Jonathan Demme', 1991, 'Crime', 8.6, 'A young FBI cadet must receive the help of an incarcerated and manipulative cannibal killer to help catch another serial killer.', 'https://image.tmdb.org/t/p/w500/uS9m8OBk1A8eM9I042bx8XXpqAq.jpg'),
    
    ('It''s a Wonderful Life', 'Frank Capra', 1946, 'Drama', 8.6, 'An angel is sent from Heaven to help a desperately frustrated businessman by showing him what life would have been like if he had never existed.', 'https://image.tmdb.org/t/p/w500/bSqt9rhDZx1Q7UZ86dBPKdNomp2.jpg'),
    
    ('Saving Private Ryan', 'Steven Spielberg', 1998, 'Drama', 8.6, 'Following the Normandy Landings, a group of U.S. soldiers go behind enemy lines to retrieve a paratrooper whose brothers have been killed in action.', 'https://image.tmdb.org/t/p/w500/uqx37cS8cpHg8U35f9U5IBlrCV3.jpg'),
    
    ('Spirited Away', 'Hayao Miyazaki', 2001, 'Animation', 8.6, 'During her family''s move to the suburbs, a sullen 10-year-old girl wanders into a world ruled by gods, witches, and spirits, where humans are changed into beasts.', 'https://image.tmdb.org/t/p/w500/39wmItIWsg5sZMyRUHLkWBcuVCM.jpg'),
    
    ('City of God', 'Fernando Meirelles', 2002, 'Crime', 8.6, 'In the slums of Rio, two kids'' paths diverge as one struggles to become a photographer and the other a kingpin.', 'https://image.tmdb.org/t/p/w500/k7eYdcdLQWYj5wALCOJ9PErcsKO.jpg'),
    
    ('Interstellar', 'Christopher Nolan', 2014, 'Adventure', 8.6, 'A team of explorers travel through a wormhole in space in an attempt to ensure humanity''s survival on a dying Earth.', 'https://image.tmdb.org/t/p/w500/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg'),
    
    ('Life Is Beautiful', 'Roberto Benigni', 1997, 'Comedy', 8.6, 'When an open-minded Jewish waiter and his son are imprisoned in a concentration camp, he uses a perfect mixture of will, humor and imagination to protect his son.', 'https://image.tmdb.org/t/p/w500/mfnkSeeVOBcGyYandeWnyJhOPEN.jpg'),
    
    ('The Green Mile', 'Frank Darabont', 1999, 'Crime', 8.6, 'The lives of guards on Death Row are affected by one of their charges: a black man accused of child murder and rape, yet who has a mysterious gift.', 'https://image.tmdb.org/t/p/w500/velWPhVMQeQKcxggNEU8YmIo52R.jpg'),
    
    ('Star Wars: A New Hope', 'George Lucas', 1977, 'Adventure', 8.6, 'Luke Skywalker joins forces with a Jedi Knight, a cocky pilot, a Wookiee and two droids to save the galaxy from the Empire''s world-destroying battle station.', 'https://image.tmdb.org/t/p/w500/6FfCtAuVAW8XJjZ7eWeLibRLWTw.jpg'),
    
    ('Terminator 2: Judgment Day', 'James Cameron', 1991, 'Action', 8.5, 'A cyborg, identical to the one who failed to kill Sarah Connor, must now protect her ten year old son, John Connor, from a more advanced and powerful cyborg.', 'https://image.tmdb.org/t/p/w500/5M0j0B18abtBI5gi2RhfjjurTqb.jpg'),
    
    ('Back to the Future', 'Robert Zemeckis', 1985, 'Adventure', 8.5, 'Marty McFly, a 17-year-old high school student, is accidentally sent thirty years into the past in a time-traveling DeLorean invented by his close friend.', 'https://image.tmdb.org/t/p/w500/fNOH9f1aA7XRTzl1sAOx9iF553Q.jpg'),
    
    ('The Pianist', 'Roman Polanski', 2002, 'Biography', 8.5, 'A Polish Jewish musician struggles to survive the destruction of the Warsaw ghetto of World War II.', 'https://image.tmdb.org/t/p/w500/2hFvxCCWrTmCwdgnwQziSDCQE2C.jpg'),
    
    ('Gladiator', 'Ridley Scott', 2000, 'Action', 8.5, 'A former Roman General sets out to exact vengeance against the corrupt emperor who murdered his family and sent him into slavery.', 'https://image.tmdb.org/t/p/w500/ty8TGRuvJLPUmAR1H1nRIsgwvim.jpg'),
    
    ('Psycho', 'Alfred Hitchcock', 1960, 'Horror', 8.5, 'A Phoenix secretary embezzles $40,000 from her employer''s client, goes on the run, and checks into a remote motel run by a young man under the domination of his mother.', 'https://image.tmdb.org/t/p/w500/yz4QVqPx3h1hD1DfqqQkCq3rmxW.jpg'),
    
    ('The Lion King', 'Roger Allers', 1994, 'Animation', 8.5, 'Lion prince Simba and his father are targeted by his bitter uncle, who wants to ascend the throne himself.', 'https://image.tmdb.org/t/p/w500/sKCr78MXSLixwmZ8DyJLrpMsd15.jpg'),
    
    ('The Departed', 'Martin Scorsese', 2006, 'Crime', 8.5, 'An undercover cop and a police informant play a dangerous game of cat and mouse with a crime boss.', 'https://image.tmdb.org/t/p/w500/nT97ifVT2J1yMQmeq20Qblg61T.jpg'),
    
    ('Whiplash', 'Damien Chazelle', 2014, 'Drama', 8.5, 'A promising young drummer enrolls at a cut-throat music conservatory where his dreams of greatness are mentored by an instructor who will stop at nothing to realize a student''s potential.', 'https://image.tmdb.org/t/p/w500/7fn624j5lj3xTme2SgiLCeuedmO.jpg'),
    
    ('The Prestige', 'Christopher Nolan', 2006, 'Drama', 8.5, 'After a tragic accident, two stage magicians engage in a battle to create the ultimate illusion while sacrificing everything they have to outwit each other.', 'https://image.tmdb.org/t/p/w500/tRNlZbgNCNOpLpbPEz5L8G8A0JN.jpg'),
    
    ('Casablanca', 'Michael Curtiz', 1942, 'Drama', 8.5, 'A cynical expatriate American cafe owner struggles to decide whether or not to help his former lover and her fugitive husband escape the Nazis in French Morocco.', 'https://image.tmdb.org/t/p/w500/5K7cOHoay2mZusSLezBOY0Qxh8a.jpg'),
    
    ('Parasite', 'Bong Joon-ho', 2019, 'Thriller', 8.5, 'Greed and class discrimination threaten the newly formed symbiotic relationship between the wealthy Park family and the destitute Kim clan.', 'https://image.tmdb.org/t/p/w500/7IiTTgloJzvGI1TAYymCfbfl3vT.jpg'),
    
    ('Alien', 'Ridley Scott', 1979, 'Horror', 8.4, 'After a space merchant vessel receives an unknown transmission as a distress call, one of the crew is attacked by a mysterious life form and they soon realize that its life cycle has merely begun.', 'https://image.tmdb.org/t/p/w500/vfrQk5IPloGg1v9Rzbh2Eg3VGyM.jpg');

-- Update the sequence to continue from the last inserted ID
SELECT setval('movies_id_seq', (SELECT MAX(id) FROM movies));

-- Display summary of inserted data
SELECT 
    COUNT(*) as total_movies,
    MIN(year) as earliest_year,
    MAX(year) as latest_year,
    ROUND(AVG(rating), 2) as average_rating,
    COUNT(DISTINCT genre) as unique_genres,
    COUNT(DISTINCT director) as unique_directors
FROM movies;