<template>
    <div class="global-container">
        <router-view/>
    </div>
</template>

<style>
.global-container {
    width: 100%;
    height: 100%;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0);
  border-radius: 0;
}

::-webkit-scrollbar {
  -webkit-appearance: none;
  width: 0;
  height: 0;
}

::-webkit-scrollbar-thumb {
  cursor: pointer;
  border-radius: 0;
  background: rgba(0, 0, 0, 0);
  transition: color 0.2s ease;
}

</style>