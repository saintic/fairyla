<template>
    <NPlayer
        :options="{
            src: fairy.src,
            videoProps: { preload: 'metadata', title: fairy.desc },
            seekStep: 5
        }"
        :set="set"
    />
</template>

<script setup>
import { defineProps } from 'vue'

defineProps({
    fairy: Object,
    set: Function
})
</script>
