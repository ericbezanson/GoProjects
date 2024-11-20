const { Sequelize } = require('sequelize');

// Create a Sequelize instance
const sequelize = new Sequelize('ebezanson', 'ebezanson', '', {
  host: '/tmp', // Use socket path
  dialect: 'postgres',
  port: 5432, // Default PostgreSQL port
});

(async () => {
  try {
    await sequelize.authenticate();
    console.log('Database connection has been established successfully.');
  } catch (error) {
    console.error('Unable to connect to the database:', error);
  }
})();

const User = require('./models/User');

// Sync database
(async () => {
  try {
    await sequelize.sync({ force: false }); // Use { force: true } for development to reset tables
    console.log('Database synchronized.');
  } catch (error) {
    console.error('Error synchronizing the database:', error);
  }
})();




module.exports = sequelize;

