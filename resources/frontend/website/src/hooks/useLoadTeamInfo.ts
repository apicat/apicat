import { onMounted } from 'vue'
import { useTeamStore } from '@/stores/team'

export const useLoadTeamInfo = (cb?: any) => {
    const teamStore = useTeamStore()

    onMounted(async () => {
        try {
            await teamStore.getTeamInfoDetail()
        } catch (e) {
            console.error('get team info error:', e)
        } finally {
            cb && cb()
        }
    })
}
