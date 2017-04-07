var TwitterPackage = require('twitter');

var creds = {
  consumer_key: '',
  consumer_secret: '',
  access_token_key: '',
  access_token_secret: ''
}
var Twitter = new TwitterPackage(creds);

var novel = [
  "It was a dark and stormy night.",
  "Suddenly, a shot rang out! A door slammed. The maid screamed.",
  "Suddenly, a pirate ship appeared on the horizon!",
  "While millions of people were starving, the king lived in luxury.",
  "Meanwhile, on a small farm in Kansas, a boy was growing up.",
  "A light snow was falling, and the little girl with the tattered shawl had not sold a violet all day.",
  "At that very moment, a young intern at City Hospital was making an important discovery.",
  "The mysterious patient in Room 213 had finally awakened. She moaned softly.",
  "Could it be that she was the sister of the boy in Kansas who loved the girl with the tattered shawl...","...who was the daughter of the maid who had escaped from the pirates?",
  "The intern frowned.",
  "\"Stampede!\" the foreman shouted, and forty thousand head of cattle thundered down on the tiny camp.", "The two men rolled on the ground grappling beneath the murderous hooves.",
  "A left and a right. A left. Another left and right. An uppercut to the jaw.",
  "The fight was over. And so the ranch was saved.",
  "The young intern sat by himself in one corner of the coffee shop.",
  "He had learned about medicine, but more importantly, he had learned something about life."
];

var index=0;

function doWork() {

  // roll a die and determine if Snoopy should tweet.
  var die = Math.floor((Math.random() * 3)); 

  if (die == 0)
    sendTweet();

}
function sendTweet() {

  // Pick a next phrase from the novel and tweet it 
  //console.log(novel[index]);

  Twitter.post('statuses/update', {status: novel[index++]}, function(error, tweet, response){
  // tweet
  }); 

  // If end of novel, reset index to beginning
  if (index == novel.length)
    index=0;

}

// send tweet on startup and set timer
sendTweet(); 
//setInterval(doWork, 10000);
setInterval(sendTweet, 2147483640);