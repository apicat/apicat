<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

use App\Repositories\ApiDoc\MockPathRepository;

class HttpApiResponseParameter extends Node
{
    protected $tagName = [];

    protected $types = ['', 'Int', 'Float', 'String', 'Array', 'Object', 'Boolean', 'File'];

    protected $projectId = '';

    protected $docId = '';

    /**
     * 是否需要project_id和doc_id
     *
     * @return boolean
     */
    public function wantIds()
    {
        return true;
    }

    public function setIds(int $projectId, int $docId)
    {
        $this->docId = $docId;
        $this->projectId = $projectId;
    }

    public function text()
    {
        if (!isset($this->node->attrs)) {
            return '';
        }

        $htmls = '';

        if (isset($this->node->attrs->response_header) and $this->node->attrs->response_header->params) {
            if (!isset($this->node->attrs->response_body) or !$this->node->attrs->response_body->params) {
                $htmls .= $this->paramTableHtml($this->node->attrs->response_header->title, $this->node->attrs->response_header->params, true);
            } else {
                $htmls .= $this->paramTableHtml($this->node->attrs->response_header->title, $this->node->attrs->response_header->params);
            }
        }

        if (isset($this->node->attrs->response_body) and $this->node->attrs->response_body->params) {
            $htmls .= $this->paramTableHtml($this->node->attrs->response_body->title, $this->node->attrs->response_body->params, true);
        }

        return $htmls;
    }

    protected function paramTableHtml($title, $params, $mock = false)
    {
        $htmls = '<div class="collapse-title">';
        $htmls .= '<h3><span class="response_body_title"><i class="iconfont iconIconCaretDown"></i>' . $title . '</span></h3>';
        if ($mock) {
            $htmls .= $this->mockUrlHtml();
        }

        $htmls .= '</div>';
        $htmls .= '<div class="ac-param-table">';
        $htmls .= '<div class="ac-param-table--border-line"></div>';
        $htmls .= '<table>';
        $htmls .= '<colgroup><col width="32%"/><col width="10%"/><col width="7%"/><col width="15%"/><col width=""/><col width="100px"/></colgroup>';
        $htmls .= '<thead><tr>';
        $htmls .= '<th>参数名称</th>';
        $htmls .= '<th>参数类型</th>';
        $htmls .= '<th class="text-center">必传</th>';
        $htmls .= '<th>默认值</th>';
        $htmls .= '<th>参数说明</th>';
        $htmls .= '<th>Mock</th>';
        $htmls .= '</tr></thead>';
        $htmls .= '<tbody>';

        foreach ($params as $index => $param) {
            $dataId = uniqid() . '-' . $index;

            $htmls .= '<tr data-id="' . $dataId . '">';

            if ($param->sub_params) {
                $htmls .= '<td>';
                $htmls .= '<i class="editor-font editor-arrow-right expand" data-id="' . $dataId . '"></i>';
                $htmls .= '<span class="copy_text">' . $param->name . '</span>';
                $htmls .= '</td>';
            } else {
                $htmls .= '<td><span class="copy_text">' . $param->name . '</span></td>';
            }
            
            $htmls .= '<td>' . $this->types[$param->type] . '</td>';
            $htmls .= '<td class="text-center">' . ($param->is_must ? '是' : '否') . '</td>';
            $htmls .= '<td>' . ($param->default_value ? $param->default_value : '') . '</td>';
            $htmls .= '<td>' . ($param->description ? $param->description : '') . '</td>';
            $htmls .= '<td><div class="mock" title="' . ((isset($param->mock_rule) and $param->mock_rule) ? $param->mock_rule : '') . '">' . ((isset($param->mock_rule) and $param->mock_rule) ? $param->mock_rule : '') . '</div></td>';
            $htmls .= '</tr>';

            if ($param->sub_params) {
                $htmls .= $this->childParamHtml($param->sub_params, $dataId);
            }
        }

        $htmls .= '</tbody></table></div>';
        
        return $htmls;
    }

    protected function childParamHtml($params, $dataPid)
    {
        $htmls = '';
        $padding = 25 + (substr_count($dataPid, '-') - 1) * 15;

        foreach ($params as $index => $param) {
            $dataId = $dataPid . '-' . $index;

            $htmls .= '<tr data-id="' . $dataId . '" data-pid="' . $dataPid . '">';

            if ($param->sub_params) {
                $htmls .= '<td style="padding-left: ' . $padding . 'px">';
                $htmls .= '<i class="editor-font editor-arrow-right expand" data-id="' . $dataId . '"></i>';
                $htmls .= '<span class="copy_text">' . $param->name . '</span>';
                $htmls .= '</td>';
            } else {
                $htmls .= '<td style="padding-left: ' . $padding . 'px"><span class="copy_text">' . $param->name . '</span></td>';
            }
            
            $htmls .= '<td>' . $this->types[$param->type] . '</td>';
            $htmls .= '<td class="text-center">' . ($param->is_must ? '是' : '否') . '</td>';
            $htmls .= '<td>' . ($param->default_value ? $param->default_value : '') . '</td>';
            $htmls .= '<td>' . ($param->description ? $param->description : '') . '</td>';
            $htmls .= '<td><div class="mock" title="' . ((isset($param->mock_rule) and $param->mock_rule) ? $param->mock_rule : '') . '">' . ((isset($param->mock_rule) and $param->mock_rule) ? $param->mock_rule : '') . '</div></td>';
            $htmls .= '</tr>';

            if ($param->sub_params) {
                $htmls .= $this->childParamHtml($param->sub_params, $dataId);
            }
        }

        return $htmls;
    }

    protected function mockUrlHtml()
    {
        $htmls = '';

        if (!$this->projectId or !$this->docId) {
            return $htmls;
        }

        if (!$record = MockPathRepository::getByDocId($this->projectId, $this->docId)) {
            return $htmls;
        }

        $url = rtrim(env('APP_MOCK_URL'), '/') . '/' . $this->projectId . '/' . ltrim($record->path, '/');

        $htmls = '<div class="http-url">';
        $htmls .= '<div class="http-url--url mock-tag">Mock</div>';
        $htmls .= '<div class="http-url--path mock-url"><a target="_blank" href="' . $url . '">' . $url . '</a></div>';
        $htmls .= '</div>';
        return $htmls;
    }
}
