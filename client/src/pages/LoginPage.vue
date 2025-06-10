<script setup lang="ts">
import { useRouter } from "vue-router"
import { repository } from "@/api"
import { ref } from "vue"

const router = useRouter()
const loginSessionRepository = repository("login-sessions")
loginSessionRepository.fetchAll()
const email = ref("")
const password = ref("")
const login = async () => {
  const r = await loginSessionRepository.create({
    email: email.value,
    password: password.value,
  })

  console.log(r)

  router.replace({
    name: "home",
  })
}
</script>

<template>
  <div class="w-full h-full items-center justify-center flex flex-col gap-4">
    <h1>Login</h1>
    <form
      method="POST"
      @submit.prevent="login"
      class="flex flex-col max-w-2xl gap-4"
    >
      <label for="email">Email</label>
      <input type="text" v-model="email" name="email" id="email" required />

      <label for="password">Password</label>
      <input
        type="password"
        v-model="password"
        name="password"
        id="password"
        required
      />

      <button
        class="bg-green-400 text-white p-4 hover:rounded hover:font-bold hover:bg-green-600 transition-colors"
        type="submit"
      >
        Login
      </button>
    </form>
    <router-link
      class="p-2 hover:rounded hover:font-bold hover:bg-yellow-200 transition-colors"
      :to="{
        name: 'register',
      }"
      >Register
    </router-link>
    <router-link
      class="p-2 hover:rounded hover:font-bold hover:bg-yellow-200 transition-colors"
      :to="{
        name: 'forgot-password',
      }"
      >Forgot Password
    </router-link>
  </div>
</template>
