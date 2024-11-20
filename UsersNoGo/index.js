const express = require('express');
const bodyParser = require('body-parser');
const db = require('./db');

const app = express();
const PORT = 3000;

// Middleware
app.use(bodyParser.json());

// POST /users - Upsert a user
app.post('/users', async (req, res) => {
  const { uid, name, email } = req.body;

  if (!uid || !name || !email) {
    return res.status(400).json({ error: 'uid, name, and email are required' });
  }

  try {
    const query = `
      INSERT INTO users (uid, name, email)
      VALUES ($1, $2, $3)
      ON CONFLICT (uid) 
      DO UPDATE SET name = $2, email = $3
      RETURNING *;
    `;
    const values = [uid, name, email];
    const { rows } = await db.query(query, values);

    res.status(200).json(rows[0]);
  }
