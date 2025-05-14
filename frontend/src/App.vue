<script setup>
import { ref } from 'vue'
import Homepage from './components/Homepage.vue'

const answer = ref('')
</script>

<template>
  <body>
    <main>
    <Homepage />
    </main>
  </body>
</template>

<style>
  #app{
    background-image: url("/frontend/src/components/icons/background.webp"), url("https://cdn.pixabay.com/photo/2021/02/22/05/00/paper-cut-6038812_1280.jpg");
    background-size: cover;
    min-width: 100%;
  }
</style>
