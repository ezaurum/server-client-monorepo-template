<template>
  <div class="w-full h-full items-center justify-center flex flex-col gap-4">
    <h1>Forgot Password</h1>
    <form method="POST" class="flex flex-col gap-4">
      <label for="email">Email</label>
      <input type="email" name="email" id="email" required />

      <button
        class="bg-green-400 text-white p-4 hover:rounded hover:font-bold hover:bg-green-600 transition-colors"
        type="submit"
      >
        Send Reset Email
      </button>
    </form>
  </div>
</template>
