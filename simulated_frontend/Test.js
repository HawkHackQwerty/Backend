const axios = require("axios");
const FormData = require("form-data");
const fs = require("fs");
const path = require("path");
const { createClient } = require("@supabase/supabase-js");

const apiBaseUrl = "http://localhost:8080"; // Adjust based on your Go server's configuration
const supabaseUrl = "https://zdgvghjjvbphcovfayyv.supabase.co";
const supabaseAnonKey =
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InpkZ3ZnaGpqdmJwaGNvdmZheXl2Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3MTYwNTk4NTMsImV4cCI6MjAzMTYzNTg1M30.fbYoUIyJkEhMImMmEgSWAdsLWAqE1-F31s2URGPqtkQ";

const supabase = createClient(supabaseUrl, supabaseAnonKey);

async function authenticate() {
  const { user, session, error } = await supabase.auth.signInWithPassword({
    email: "test@gmail.com",
    password: "Test123",
  });

  if (error)
    throw new Error("Failed to authenticate with Supabase:", error.message);

  const res = await supabase.auth.getSession();
  const sessionData = res?.data.session;

  const a = sessionData.user.id;
  const b = sessionData.access_token;

  return {
    userId: a,
    authToken: b,
  };
}

async function uploadFile(endpoint, filePath, fileType, headers) {
  const formData = new FormData();
  formData.append(
    "file",
    fs.createReadStream(filePath),
    path.basename(filePath),
  );

  try {
    const response = await axios.post(`${apiBaseUrl}/${endpoint}`, formData, {
      headers: {
        ...formData.getHeaders(),
        ...headers,
      },
      responseType: "json",
    });
    console.log(`${endpoint} Response:`, response.data);
  } catch (error) {
    console.error(`${endpoint} Error:`, error.response.data);
  }
}

async function uploadJobInfo(headers) {
  try {
    const response = await axios.post(
      `${apiBaseUrl}/uploadJob`,
      {
        stringOne: "Hello",
        stringTwo: "World",
      },
      {
        headers: {
          "Content-Type": "application/json",
          ...headers,
        },
      },
    );
    console.log(`uploadJob Response:`, response.data);
  } catch (error) {
    console.error(`uploadJob Error:`, error.response.data);
  }
}

async function runTests() {
  const { userId, authToken } = await authenticate();

  console.log(userId);
  console.log(authToken);
  const headers = {
    Authorization: `Bearer ${authToken}`,
    "X-User-ID": userId,
  };

  await uploadFile(
    "processResume",
    "./Resuume.pdf",
    "application/pdf",
    headers,
  );
  await uploadFile("processCover", "./Cover.pdf", "application/pdf", headers);
  await uploadFile("processVideo", "./video.mp4", "video/mp4", headers);
  await uploadJobInfo(headers);
}

runTests();
