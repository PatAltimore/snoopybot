const { app } = require('@azure/functions');
const { doWork } = require('../bot');

app.timer('snoopyTrigger', {
  schedule: '0 0 17 * * *', // Daily at 5:00 PM UTC
  handler: async (myTimer, context) => {
    context.log('Snoopy is writing...');
    await doWork(context);
    context.log('Snoopy is done.');
  },
});
