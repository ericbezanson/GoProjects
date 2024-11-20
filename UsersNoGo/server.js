const express = require('express');
const bodyParser = require('body-parser');
const path = require('path');
const User = require('./models/User'); // Assuming the User model is in ./models/User.js

const app = express();
const port = 3000;

// Middleware
app.use(bodyParser.json());

// Serve the HTML form at the root URL
app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, 'index.html')); // Adjust path if your file structure is different
});

// POST /users - Create a new user
app.post('/users', async (req, res) => {
  const { name, email } = req.body;

  try {
    const newUser = await User.create({ name, email });
    console.log("NEW USERS", newUser)
    res.status(201).json(newUser); // Return the created user as JSON
  } catch (error) {
    console.error('Error creating user:', error);
    res.status(500).json({ error: 'An error occurred while creating the user' });
  }
});

// Start the server
app.listen(port, () => {
  console.log(`Server running on http://localhost:${port}`);
});
