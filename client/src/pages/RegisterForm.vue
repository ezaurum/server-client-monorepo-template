<template>
  <div class="w-full h-full items-center justify-center flex flex-col gap-4">
    <h1>Register</h1>
    <form method="POST" class="flex flex-col max-w-2xl gap-4">
      <label for="email">Email</label>
      <input type="email" name="email" id="email" required />

      <label for="password">Password</label>
      <input type="password" name="password" id="password" required />

      <button
        class="bg-green-400 text-white p-4 hover:rounded hover:font-bold hover:bg-green-600 transition-colors"
        type="submit"
      >
        Register
      </button>
    </form>

    <router-link
      :to="{
        name: 'login',
      }"
      class="p-2 hover:rounded hover:font-bold hover:bg-yellow-200 transition-colors"
    >
      Login
    </router-link>
  </div>
</template>
