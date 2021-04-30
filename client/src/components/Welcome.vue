<template>
    <div>
        <p>欢迎~</p>
    </div>
</template>

<script setup>
import { defineProps, reactive } from 'vue'

defineProps({
    msg: String
})

const state = reactive({ count: 0 })
</script>

<style scoped>
a {
    color: #42b983;
}
</style>
