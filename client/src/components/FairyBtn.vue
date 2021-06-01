<template>
    <el-badge
        class="func-btn"
        :is-dot="btn.badge.dot"
        :value="btn.badge.value"
        :type="btn.badge.type"
        v-if="btn.badge"
    >
        <el-button
            size="mini"
            :type="btn.type"
            :plain="btn.plain"
            :disabled="btn.disabled"
            :icon="btn.icon"
            @click="btn.click"
            >{{ btn.name }}</el-button
        >
    </el-badge>
    <el-button
        v-else
        size="mini"
        :type="btn.type"
        :plain="btn.plain"
        :disabled="btn.disabled"
        :icon="btn.icon"
        @click="btn.click"
        >{{ btn.name }}</el-button
    >
</template>

<script setup>
import { defineProps } from 'vue'
import { isObject } from '@/libs/util.js'

defineProps({
    btn: {
        type: Object,
        validator: (value) => {
            if (!isObject(value)) return false
            if (!value.hasOwnProperty('name')) return false
            return true
        }
    }
})
</script>

<style scoped>
.func-btn {
    margin-right: 10px;
}
</style>