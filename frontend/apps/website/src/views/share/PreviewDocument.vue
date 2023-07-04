<template>
  <div :class="[ns.b(), { 'h-20vh': !httpDoc }]" v-loading="isLoading">
    <div class="ac-editor mt-10px" v-if="httpDoc">
      <RequestMethodRaw class="mb-10px" :doc="httpDoc" :urls="urlServers" />
      <RequestParamRaw class="mb-10px" :doc="httpDoc" :definitions="definitions" />
      <ResponseParamTabsRaw :doc="httpDoc" :definitions="definitions" />
    </div>
  </div>
</template>
<script setup lang="ts">
import { HttpDocument } from '@/typings'
import { useNamespace } from '@/hooks/useNamespace'
import ResponseParamTabsRaw from '@/components/ResponseParamTabsRaw.vue'
import { storeToRefs } from 'pinia'
import { getCollectionDetail } from '@/api/collection'
import useDefinitionStore from '@/store/definition'
import useDefinitionResponseStore from '@/store/definitionResponse'
import uesProjectStore from '@/store/project'
import uesGlobalParametersStore from '@/store/globalParameters'

const projectStore = uesProjectStore()
const definitionStore = useDefinitionStore()
const globalParametersStore = uesGlobalParametersStore()
const definitionResponseStore = useDefinitionResponseStore()

const [isLoading, getCollectionDetailApi] = getCollectionDetail()
const { urlServers } = storeToRefs(projectStore)
const { definitions } = storeToRefs(definitionStore)

const ns = useNamespace('document')
const httpDoc: Ref<HttpDocument | null> = ref(null)

const getDetail = async (docId: string) => {}
</script>
