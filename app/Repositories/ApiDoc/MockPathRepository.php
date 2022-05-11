<?php

namespace App\Repositories\ApiDoc;

use App\Models\MockPath;

class MockPathRepository
{
    public static $methods = [
        '',
        'get',
        'post',
        'put',
        'patch',
        'delete',
        'option'
    ];

    /**
     * 通过文档id获取记录
     * @param int $projectId 项目id
     * @param int $docId 文档id
     * @return MockPath
     */
    public static function getByDocId($projectId, $docId)
    {
        return MockPath::where([
            ['project_id', $projectId],
            ['doc_id', $docId]
        ])->first();
    }

    /**
     * 更新api mock的路径
     * @param int $projectId 项目id
     * @param int $docId 文档id
     * @param string $path api路径
     * @param int $method 请求方法
     * @return boolean
     */
    public static function updatePath($projectId, $docId, $path, $method)
    {
        if (!$path) {
            return self::del($projectId, $docId);
        }

        $record = MockPath::withTrashed()->where([
            ['project_id', $projectId],
            ['doc_id', $docId]
        ])->first();

        $path = '/' . ltrim($path, '/');

        if ($record) {
            $record->path = $path;
            $record->method = self::$methods[$method];
            $record->deleted_at = null;
            return $record->save();
        }

        return (bool)MockPath::create([
            'project_id' => $projectId,
            'doc_id' => $docId,
            'path' => $path,
            'format' => 'json',
            'method' => self::$methods[$method]
        ]);
    }

    /**
     * 更新api mock的数据格式
     * @param int $projectId 项目id
     * @param int $docId 文档id
     * @param string $format 数据格式
     * @return boolean
     */
    public static function updateFormat($projectId, $docId, $format)
    {
        return (bool)MockPath::withTrashed()->where([
            ['project_id', $projectId],
            ['doc_id', $docId]
        ])->update(['format' => $format, 'deleted_at' => null]);
    }

    /**
     * 删除api mock记录
     * @param int $projectId 项目id
     * @param int $docId 文档id
     * @return boolean
     */
    public static function del($projectId, $docId)
    {
        return (bool)MockPath::where([
            ['project_id', $projectId],
            ['doc_id', $docId]
        ])->delete();
    }

    /**
     * 恢复api mock记录
     * @param int $projectId 项目id
     * @param int $docId 文档id
     * @return boolean
     */
    public static function restore($projectId, $docId)
    {
        $record = MockPath::withTrashed()->where([
            ['project_id', $projectId],
            ['doc_id', $docId]
        ])->first();
        if (!$record) {
            return false;
        }

        $record->restore();

        return true;
    }
}