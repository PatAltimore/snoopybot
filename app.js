var Twit = require('twit');
var redis = require('redis');
var express = require('express');
var app = express();

// Twitter & redis cache credentials
var config = require('./config.js');
var Twitter = new Twit(config);
var cache = redis.createClient(6380, config.cache.redis_servername, {auth_pass: config.cache.redis_auth_pass, tls: {servername: config.cache.redis_servername}});

// Snoopy's novels 
var novel = [
  'It was a dark and stormy night.',
  'Suddenly, a shot rang out! A door slammed. The maid screamed.',
  'Suddenly, a pirate ship appeared on the horizon!',
  'While millions of people were starving, the king lived in luxury.',
  'Meanwhile, on a small farm in Kansas, a boy was growing up.',
  'Part II - A light snow was falling, and the little girl with the tattered shawl had not sold a violet all day.',
  'At that very moment, a young intern at City Hospital was making an important discovery.',
  'The mysterious patient in Room 213 had finally awakened. She moaned softly.',
  'Could it be that she was the sister of the boy in Kansas who loved the girl with the tattered shawl...','...who was the daughter of the maid who had escaped from the pirates?',
  'The intern frowned.',
  '"Stampede!" the foreman shouted, and forty thousand head of cattle thundered down on the tiny camp.', 'The two men rolled on the ground grappling beneath the murderous hooves.',
  'A left and a right. A left. Another left and right. An uppercut to the jaw.',
  'The fight was over. And so the ranch was saved.',
  'The young intern sat by himself in one corner of the coffee shop.',
  'He had learned about medicine, but more importantly, he had learned something about life. THE END',
  '"A Love Story" by Erich Beagle: "I love you," she said, and together they laughed.',
  'Then one day she said, "I hate you," and they cried. But not together.',
  '"What happened to the love that we said would never die?" she asked. "It died," he said.',
  'The first time he saw her she was playing tennis. The last time he saw her she was playing tennis.',
  '"Ours was a Love set," he said, "but we double faulted." "You always talked a better game than you played," she said.',
  'Though her husband often went on business trips, she hated to be left alone.',
  '"I\'ve solved our problem," he said. "I\'ve bought you a St. Bernard. Its name is Great Reluctance.',
  'Now, when I go away, you shall know that I am leaving you with Great Reluctance!" She hit him with a waffle iron.'
];

// Miscellanous quotes and writings
var misc = [
  "Here's the world-famous author waiting for word from his publisher...",
  "Sometimes, when you are a great writer, the words come so fast you can hardly put them down on paper...",
  "My plot is thickening.",
  "Stop raining on my novel!",
  "This twist in the plot will baffle my readers.",
  "In Part Two, I'll tie all of this together.",
  "I may have written myself into a corner.",
  "For the first time in my life I know how Leo must have felt. Leo Tolstoy, that is!",
  "Here's the world-famous novelist walking to the mailbox to send his manuscript away.",
  "I'm expecting word from my publisher.",
  "I can see the headlines now, \"Author bites book reviewer on the leg!\"",
  "Here I am on the way to my first autograph party. I hate it when the line doesn't extend clear around the block.",
  "Autograph parties are terrible when nobody shows up.",
  "What I need is a unique signature. 🐾",
  "It was a dark and stormy night. Suddenly, a kiss rang out!",
  "It was a dark and stormy night. Suddenly, a vote rang out!",
  "It was a dark and stormy night. Suddenly, a turkey rang out!",
  "Once upon a time...It was a dark and stormy night.",
  "It was a dark and stormy night. Suddenly, out of the mist a spooky figure appeared. How spooky was he? Spoooooooky!",
  "I have the perfect title... \"Has It Ever Occurred to You That You Might Be Wrong?\""
  ];

function doWork() {

  var tweet, index;

  // roll a die and determine what Snoopy should tweet.
  var die = Math.floor((Math.random() * 2)); 

  switch(die) {

    // Tweet next sentence in the novel
    case 0:
      // Retrieve last novel index from cache
      cache.get('sbNovelIndex',  function(err, index) {

        // Increment index and store in cache
        if (++index == novel.length)
          index=0;

        cache.set('sbNovelIndex', index, redis.print);
        
        // Tweet sentence
        console.log('Tweeting novel index:' + index);
        tweet = sendTweet(novel[index]); 
      });
      break;

    // Tweet random quote
    case 1:
      var range = misc.length-1;
      var phrase = Math.floor((Math.random() * range) + 1);
      
      console.log('Tweeting miscellaneous quote index:' + phrase);
      tweet = sendTweet(misc[phrase]);
      break;
  }
  return tweet;
}

// Call Twitter API to update status (Tweet)
function sendTweet(tweet) {

  Twitter.post('statuses/update', { status: tweet}, function(error, tweetResponse, response){
    if(error){
      tweetPat(error);
      console.log(error);
    }
  }); 

  // Save last tweet to cache
  cache.set('sbLastAction', tweet, redis.print);

  console.log(tweet);
  return tweet;
}

// Send direct message to Pat if error
function tweetPat(error)
{
  // send error to Pat
  for (var i = 0; i < error.length; i++) 
  {
    lastError = error[i].code + ': ' + error[i].message;
    Twitter.post('direct_messages/new', { screen_name: '@PatAltimore', text: lastError }, function(error, tweetResponse, response) {
      if(error){
        console.log(error);
      }
    });
    cache.set('sbLastAction', error[i].message, redis.print);
  }
}

// Return "Hello" message if request to website
app.get('/', function (req, res) {
  res.send('<h1>Hello, I\'m the Snoopy twitter bot.</h1>');
})


// Return last action status for /status resource
app.get('/status', function (req, res) {

  cache.get('sbLastAction',  function(err, reply) {
        status = '<h1>Snoopy twitter bot status</h1>';
        status += 'Last action: ' + reply;

        res.send(status);
      });
})

// Force tweet
/* app.get('/tweet', function (req, res) {
  var tweet = '<h1>Snoopy just tweeted</h1>';
  doWork();
  res.send(tweet);
}) */

// If followed, reply back to follower and follow them back
function followed(event) {
  var name = event.source.name;
  var screenName = event.source.screen_name;

  if (screenName !== "SnoopyAtWork") // Don't follow yourself
  {
    console.log('Followed by: ' + name + ' ' + screenName); 

    // direct message back to follower 
    var tweet = 'Happiness is being followed. ❤️';

    Twitter.post('direct_messages/new', { screen_name: screenName, text: tweet }, function(error, tweetResponse, response){
      if(error){
        tweetPat(error);
        console.log(error);
      }
    }); 

    // Follow them back
    Twitter.post('friendships/create', { screen_name: screenName}, function(error, tweetResponse, response){
      if(error){
        tweetPat(error);
        console.log(error);
      }
    }); 
  }
}

// If set up user stream, listen for followers...
var stream = Twitter.stream('user');
stream.on('follow', followed);

// Listen on port
app.listen(process.env.PORT || 8080);