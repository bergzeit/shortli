<script setup>
import { ref } from 'vue'

const props = defineProps({
  msg: String
})

import { config } from "../config.js"
const textInput = ref('')
const textInputLong = ref('')
let answer = ref('')
let shortUrl = ref('')
let longUrl = ref('')
let firstPartLink = "http://localhost:8080/"
let forwardingLink = ref('')

// fetchData catches the JSON-Data from Golang.
// input must be an existing shortlink in the db. database.
async function fetchData(input) {
    try {
        // response is waiting for the answer from server (Golang).
        const response = await fetch(config.url + input);
        if (!response.ok) { //response status code must be valid.
            throw new Error('Network response was not ok');
        };
        const data = await response.json(); // data is the hole JSON
        return data;
    } catch (error) {
        console.error('Error while loading the data:', error);
        return null;
    }
}

async function fetchOriginalLink() {
  const result = await fetchData(textInput.value);
  if (result) {
    answer.value = result.originalUrl;  // answer is now the original link as string
    shortUrl.value = result.shortUrl;
    forwardingLink = firstPartLink + shortUrl.value
  } else {
    answer.value = "Link existiert nicht."
    shortUrl.value = "Shortlink existiert nicht."
  };

  resetInputField()
}

// handelSubmit sends given original link (JSON) to backend.
// Golang sends afterwards the answer (new shortlink) to the client.
async function handleSubmit() {
  const data = {
    originalUrl: textInputLong.value
  }

  const postResponse = await fetch(config.post, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
  });

  // Backend sends the result.
  const result = await postResponse.json();
  shortUrl.value = result.shortUrl;
  longUrl.value = data.originalUrl;
  answer.value = "";
  forwardingLink = firstPartLink + shortUrl.value

  resetInputField()
}

function resetInputField() {
  textInputLong.value = ""
  textInput.value = ""
}
</script>

<template>  
  <body>
    <div class="body-elements">
      <div class="title-text">
        Generate a short Link!
      </div>

      <!--POST-METHOD-->
      <form @submit.prevent="handleSubmit()">
        <input 
          v-model="textInputLong" 
          type="url"
          name="url" 
          placeholder="https://example.com" 
          required
          pattern="https?://.*"
          title="URL INVALID (http:// or https://)"
        />
        <button type="submit">
          Kurzen Link erstellen
        </button>
        <p v-if="longUrl">
          Originaler Link: <br><a class="generatedLinks" :href="`${longUrl}`"> {{ longUrl }}</a><br>
          Verkürzter Link: <br><a class="generatedLinks" :href="`${forwardingLink}`"> {{ shortUrl }}</a><br>
        </p>
      </form>
   
      <!--GET-METHOD-->
      <form @submit.prevent="fetchOriginalLink()">
        <input 
          v-model="textInput" 
          type="text" 
          name="shortlink"
          placeholder="Shortlink eingeben"
          required
        />
        <button type="submit">
          Originalen Link anzeigen
        </button>
        <p v-if="answer">
          Originaler Link: <a class="generatedLinks" :href="`${answer}`"><br> {{ answer }} </a><br>
          Verkürzter Link: <br><a class="generatedLinks" :href="`${forwardingLink}`"> {{ shortUrl }}</a><br>       
        </p>
      </form>

    </div>
  </body>
</template>

<style>
  .title-text{
    color: white;
    font-size:  50px;
  }
  button {
    background-color: rgb(255, 146, 4);
    color: white;
    border: none;
    padding: 10px;
    border-radius: 20px;
    margin: 10px;
    cursor: pointer;
    font-weight: bold;
  }
  .body-elements{
    border-style: solid;
    border-color: white;
    padding: 25px;
    background-color: rgba(67, 0, 65, 0.9);
    max-width: 80vh;
  }
  html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  }
  body{
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    background-color: transparent;
    min-height: 100vh;
  }
  div#app {
  width: 100%;
  max-width: 800px;
  padding: 20px;
  box-sizing: border-box;
  min-height: 60vh; /* oder 100vh für vollen Viewport */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  }
  input {
    height: 30px;
    width: 300px;
  }
  .generatedLinks{
    font-weight: bold;
    color: cornsilk;
  }
  a.generatedLinks {
  display: inline;
  max-width: 100%;
  white-space: normal;
  overflow-wrap: break-word;
  word-break: break-word; 
}
  a,
  p {
    font-size: 30px;
    color: rgb(255, 255, 255);
    background-color: transparent;
    text-decoration: none;
    margin: 15px ;
    margin-left: 0px;
  }
  a:hover {
    background-color: transparent;
    text-decoration:underline;
  }

</style>