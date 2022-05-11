<template>
    <div class="image-preview cursor-pointer">
        <input type="file" ref="input" class="upload-file-image" accept=".png, .jpg, .jpeg" />
        <slot></slot>
    </div>
</template>
<script>
    export default {
        mounted() {
            var fileInput = this.$refs['input']
            fileInput &&
                fileInput.addEventListener('change', (e) => {
                    var files = e.target.files
                    var reader
                    var file

                    var done = (url) => {
                        fileInput.value = ''
                        this.$emit('done', url, file)
                    }

                    if (files && files.length > 0) {
                        file = files[0]

                        if (file && file.size > 1 * 1024 * 1024) {
                            this.$Message.error('文件不能超过1MB。')
                            fileInput.value = ''
                            return
                        }

                        if (URL) {
                            done(URL.createObjectURL(file), file)
                        } else if (FileReader) {
                            reader = new FileReader()
                            reader.onload = function (e) {
                                done(reader.result, file)
                            }
                            reader.readAsDataURL(file)
                        }
                    }
                })
        },
    }
</script>

<style lang="scss">
    .image-preview {
        display: inline-block;
        position: relative;
        overflow: hidden;
        direction: ltr;
        cursor: pointer;

        .upload-file-image {
            font-size: 100%;
            opacity: 0;
            position: absolute;
            width: 100%;
            height: 100%;
            cursor: pointer;
        }
    }
</style>
