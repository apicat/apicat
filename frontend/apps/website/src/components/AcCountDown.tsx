export default defineComponent({
  props: {
    link: {
      type: String,
      default: '',
    },
    time: {
      type: Number,
      default: 0,
    },
    autoJump: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const router = useRouter()
    const countDown = ref(props.time)
    let timer: any = null

    function clearTimer() {
      timer && clearTimeout(timer)
    }

    function jump() {
      props.link && router.replace(props.link)
    }

    onMounted(() => {
      timer = setInterval(() => {
        countDown.value--
        if (countDown.value <= 0) {
          props.autoJump && jump()
          return clearTimer()
        }
      }, 1000)
    })

    onUnmounted(() => clearTimer())

    const slots = useSlots()
    return () => {
      return (<>{slots.default?.({ seconds: countDown.value })}</>)
    }
  },
})
