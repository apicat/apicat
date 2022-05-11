<?php

namespace App\Modules\EditorJsonToHtml\Marks;

class Mark
{
    protected $mark;

    protected $tagName = null;

    public function __construct($mark)
    {
        $this->mark = $mark;
    }

    public function tag()
    {
        return $this->tagName;
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

        return implode('', array_map(function($item) {
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
}
