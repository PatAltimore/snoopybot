var Twit = require('twit');
var storage = require('azure-storage');

// Twitter credentials
var config = require('./config.js');

var Twitter = new Twit(config);

// Create table
var tableSvc = storage.createTableService(config.storage.AZURE_STORAGE_ACCOUNT, config.storage.AZURE_STORAGE_ACCESS_KEY);

tableSvc.createTableIfNotExists('state', function(error, result, response) {
  if(!error){
    
    // If the table is new, add an entry for novel state
    if (result.created) {
      var entGen = storage.TableUtilities.entityGenerator;
      var state = {
        PartitionKey: entGen.String('state'),
        RowKey: entGen.String('novel'),
        index: entGen.Int32(0)
      }
      tableSvc.insertEntity('state', state, function(error, result, response) {
        if(!error) {
          console.log('Storage table and initial entry created');
        }
      })
    }
  }
})

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
  "It was a dark and stormy night. Suddenly, out of the mist a spooky figure appeared. How spooky was he? Spoooooooky!",
  "I have the perfect title... \"Has It Ever Occurred to You That You Might Be Wrong?\""
  ];

  function doWork() {

    var tweet, index=0;
  
    // roll a die and determine what Snoopy should tweet.
    var die = Math.floor((Math.random() * 2)); 
  
    switch(die) {
  
      // Tweet next sentence in the novel
      case 0:
        // Retrieve last novel index from storage table
        tableSvc.retrieveEntity('state', 'state', 'novel', function(error, state, response) {
          if(!error){
            index = state.index._;
          }
          // Increment index
          index = index + 1;
          if (index > novel.length-1)
            index=0;
  
          // Update storage table with new index
          state.index._ = index;
  
          tableSvc.replaceEntity('state', state, function(error, result, response) {
            if(!error) {
              // entity updated
            }
          })
  
          // Tweet sentence
          console.log('Tweeting novel index:' + index);
          sendTweet(novel[index]); 
        });
        break;
  
      // Tweet random quote
      case 1:
        var range = misc.length-1;
        var phrase = Math.floor((Math.random() * range) + 1);
        
        console.log('Tweeting miscellaneous quote index:' + phrase);
        sendTweet(misc[phrase]);
        break;
    }
  }
  
  // Call Twitter API to update status (Tweet)
  function sendTweet(tweet) {
  
     Twitter.post('statuses/update', { status: tweet}, function(error, tweetResponse, response){
      if(error){
        tweetPat(error);
        console.log(error);
      }
    }); 
  
    console.log(tweet);
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
    }
  }

// For webjob, tweet and exit
doWork();

// Exit after 20 seconds
setTimeout(function () {
    process.exit();
}, 20000);
