export default {
  props: {
    editor:{
      type:Object,
      default:()=>({})
    },
    isCreate:{
      type:Boolean,
      default:false
    }
  },
  data() {
    return {
      node: null,
      attrs: {}
    };
  },

  methods: {
    setNode({ node }) {
      if (!node) {
        console.error("node 不能为空");
        return;
      }
      this.node = node;
      this.attrs = { ...node.attrs };
    },

    close(){
      this.$emit('on-close')
    }
  },
};
