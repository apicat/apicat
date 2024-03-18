import { useGlobalLoading } from "./useGlobalLoading";

type AsyncLifecycleHook = () => Promise<void> | void;

export function useInitedPageWithGlobalLoading(initLogic: AsyncLifecycleHook): AsyncLifecycleHook {

    return async () => {
        if (!initLogic) {
            return
        }
        const { showGlobalLoading, hideGlobalLoading } = useGlobalLoading()
        try {
            showGlobalLoading()
            await initLogic?.()
        } finally {
            hideGlobalLoading()
        }
    }
}