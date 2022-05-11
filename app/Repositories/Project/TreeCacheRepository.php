<?php

namespace App\Repositories\Project;

use Illuminate\Support\Facades\Cache;

class TreeCacheRepository
{
    /**
     * 获取项目文档树的缓存
     * @param int $projectID 项目id
     * @return array|null
     */
    public static function get($projectID)
    {
        $cacheKey = self::cacheKey($projectID);
        return Cache::store('file')->get($cacheKey);
    }

    /**
     * 设置项目文档树的缓存
     * @param int $projectID 项目id
     * @param array $tree 树信息
     * @return boolean
     */
    public static function set($projectID, $tree)
    {
        $cacheKey = self::cacheKey($projectID);
        return Cache::store('file')->put($cacheKey, $tree, now()->addDay());
    }

    /**
     * 删除项目文档树的缓存
     * @param int $projectID 项目id
     * @return boolean
     */
    public static function remove($projectID)
    {
        $cacheKey = self::cacheKey($projectID);
        return Cache::store('file')->forget($cacheKey);
    }

    /**
     * 获取缓存key
     * @param int $projectID 项目id
     * @return string
     */
    public static function cacheKey($projectID)
    {
        return 'project_tree:' . $projectID;
    }
}
