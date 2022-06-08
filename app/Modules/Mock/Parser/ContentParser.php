<?php

namespace App\Modules\Mock\Parser;

use Spatie\ArrayToXml\ArrayToXml;

/**
 * 解析API文档内容
 */
class ContentParser
{
    /**
     * 解析API文档内容为可mock的数据格式
     *
     * @param object $content API文档内容
     * @return array
     */
    public static function make($content)
    {
        $responseHeader = $responseBody = [];

        if ($content->response_header->params) {
            foreach ($content->response_header->params as $param) {
                if (!$param->name or in_array($param->type, [4, 5, 6, 7, 8])) {
                    continue;
                }
                $responseHeader[$param->name] = ParamParser::parse($param);
            }
        }

        if ($content->response_body->params) {
            foreach ($content->response_body->params as $param) {
                if (!$param->name) {
                    continue;
                }
                $responseBody[$param->name] = ParamParser::parse($param);
            }
        }

        return [
            'response_header' => $responseHeader,
            'response_body' => $responseBody
        ];
    }

    /**
     * 生成mock数据
     *
     * @param array $content 可mock的数据
     * @param string $format 返回的mock数据格式: array,json,xml
     * @return array|string
     */
    public static function build($content, $format = 'json')
    {
        if (!$content) {
            return $format == 'array' ? [] : '';
        }

        $result = [];
        foreach ($content as $k => $v) {
            $result[$k] = Maker\MockRouter::generateData($v);
        }

        if ($format == 'json') {
            return json_encode($result);
        } elseif ($format == 'array') {
            return $result;
        } else {
            // xml待实现
            return ArrayToXml::convert($result);
        }
    }

    /**
     * 生成mock数据
     *
     * @param object $content API文档内容
     * @param string $format 返回的mock数据格式: json,xml
     * @return string
     */
    public static function mock($content, $format = 'json')
    {
        $response = self::make($content);
        $mockHeader = self::build($response['response_header'], 'array');
        $mockBody = self::build($response['response_body'], $format);
        return [
            'mock_header' => $mockHeader,
            'mock_body' => $mockBody
        ];
    }
}
