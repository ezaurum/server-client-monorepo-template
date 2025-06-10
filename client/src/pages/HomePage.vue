<template>
  <div>
    <nav>
      <router-link to="/me">My profile</router-link>
      <router-link to="logout">Logout</router-link>
      <router-link to="sheet/new">새 출석부 만들기</router-link>
    </nav>
    <div>
      <Button class="bg-white text-black" @click.prevent="test">Test</Button>
      <h1>Home</h1>
      <p>Welcome to the home page</p>
      <ol>
        <li>출석부 1</li>
        <li>출석부 2</li>
        <li>출석부 3</li>
      </ol>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { Button } from "@/components/ui/button"

const test = () => {
  console.log("click")
}
</script>
