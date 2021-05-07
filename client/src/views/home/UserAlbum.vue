<template>
    <Album :albums="albums" />
</template>

<script>
import Album from '@/components/Album.vue'

export default {
    Name: 'UserAlbum',
    components: { Album },
    data() {
        return { albums: [] }
    },
    created() {
        this.$http.get('/user/album').then((res) => {
            for (let a of res.data) {
                if (a.steady_fairy) {
                    a['fairy'] = a.steady_fairy
                } else {
                    a['fairy'] = a.latest_fairy
                }
                this.albums.push(a)
            }
        })
    }
}
</script>
