<template>
    <el-result title="404" subTitle="页面未发现">
        <template #icon>
            <el-image :src="img404"></el-image>
        </template>
    </el-result>
</template>

<script setup>
import img404 from '@/assets/img/404.png'
</script>
