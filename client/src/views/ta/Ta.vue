<template>
    <Album :albums="albums" />
</template>

<script>
import Album from '@/components/Album.vue'

export default {
    Name: 'Ta',
    components: { Album },
    data() {
        return { albums: [] }
    },
    created() {
        this.$http.get('/album').then((res) => {
            for (let uas of res.data) {
                for (let a of uas) {
                    if (a.steady_fairy) {
                        a['fairy'] = a.steady_fairy
                    } else {
                        a['fairy'] = a.latest_fairy
                    }
                    this.albums.push(a)
                }
            }
        })
    }
}
</script>
