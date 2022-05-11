export default (content, schema) => {
  let unknownTypes = [];
  if (content.type === "doc" && content.content && content.content.length) {
    content.content = Array.from(content.content).filter((node) => {
        if(schema.nodes[node.type] !== undefined){
            return true
        }
        unknownTypes.push(node.type)
      return false;
    });

    unknownTypes.length && console.error("Unknown node type:",unknownTypes.join("、"))
    return content;
  }

  return content;
};
