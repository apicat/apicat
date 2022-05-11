<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class Node
{
    /**
     * 节点对象
     *
     * @var object
     */
    protected $node;

    /**
     * 是否以打印格式渲染
     *
     * @var boolean
     */
    protected $print;

    /**
     * 标签名称
     *
     * @var string
     */
    protected $tagName = null;

    /**
     * construct
     *
     * @param object $node 节点对象
     * @param boolean $print 是否是打印格式
     */
    public function __construct($node, $print)
    {
        $this->node = $node;
        $this->print = $print;
    }

    public function tag()
    {
        return $this->tagName;
    }

    /**
     * 是否为自闭合标签
     *
     * @return boolean
     */
    public function selfClosing()
    {
        return false;
    }

    /**
     * 自定义内容
     *
     * @return string
     */
    public function text()
    {
        return '';
    }

    /**
     * 渲染html开始标签
     *
     * @return string
     */
    public function renderOpeningTag()
    {
        $tags = $this->tag();
        $tags = (array)$tags;

        if (!$tags or count($tags) < 1) {
            return '';
        }

        return implode('', array_map(function ($item) {
            if (is_string($item)) {
                return '<' . $item . '>';
            }

            $attrs = '';
            if (isset($item['attrs'])) {
                foreach ($item['attrs'] as $attribute => $value) {
                    $attrs .= ' ' . $attribute . '=' . '"' . $value . '"';
                }
            }

            return '<' . $item['tag'] . $attrs . '>';
        }, $tags));
    }

    /**
     * 渲染html结束标签
     *
     * @return string
     */
    public function renderClosingTag()
    {
        $tags = $this->tag();
        $tags = (array)$tags;
        $tags = array_reverse($tags);

        if (!$tags or count($tags) < 1) {
            return '';
        }

        return implode('', array_map(function ($item) {
            if (is_string($item)) {
                return '</' . $item . '>';
            }

            return '</' . $item['tag'] . '>';
        }, $tags));
    }

    /**
     * 是否需要project_id和doc_id
     *
     * @return boolean
     */
    public function wantIds()
    {
        return false;
    }

    /**
     * @param int $projectId 项目id
     * @param int $docId 文档id
     */
    public function setIds(int $projectId, int $docId)
    {
    }
}
