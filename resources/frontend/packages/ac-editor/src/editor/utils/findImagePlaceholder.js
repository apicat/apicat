/**
 * @param {} plugin
 * @param {*} state
 * @param {*} id
 * @returns 插件装饰集合
 */
export default function findImagePlaceholder(plugin, state, id) {
  const decos = plugin.getState(state);
  const found = decos.find(null, null, (spec) => spec.id === id);
  return found.length ? found[0].from : null;
}
