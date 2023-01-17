<?php

namespace App\Repositories\Project;

use Illuminate\Support\Facades\Cache;

class TreeCacheRepository
{
    /**
     * 获取项目文档树的缓存
     * @param int $projectID 项目id
     * @param int $iterationId 迭代id
     * @return array|null
     */
    public static function get(int $projectID, int $iterationId = 0)
    {
        if ($iterationId) {
            $iterationListCacheKey = self::iterationListCacheKey($projectID);
            $iterationCacheKey = self::iterationCacheKey($iterationId);

            if ($cache = Cache::store('file')->get($iterationListCacheKey)) {
                return isset($cache[$iterationCacheKey]) ? $cache[$iterationCacheKey] : null;
            }
            return;
        }

        $projectCacheKey = self::projectCacheKey($projectID);
        return Cache::store('file')->get($projectCacheKey);
    }

    /**
     * 设置文档树的缓存
     *
     * @param int $projectID 项目id
     * @param array $tree 树信息
     * @param int $iterationId 迭代id
     * @return boolean
     */
    public static function set(int $projectID, array $tree, int $iterationId = 0)
    {
        if ($iterationId) {
            return self::setIteration($projectID, $tree, $iterationId);
        }

        $projectCacheKey = self::projectCacheKey($projectID);
        return Cache::store('file')->put($projectCacheKey, $tree, now()->addDay());
    }

    /**
     * 设置迭代文档树的缓存
     *
     * @param int $projectID 项目id
     * @param array $tree 树信息
     * @param int $iterationId 迭代id
     * @return boolean
     */
    public static function setIteration(int $projectID, array $tree, int $iterationId)
    {
        $iterationListCacheKey = self::iterationListCacheKey($projectID);
        $iterationCacheKey = self::iterationCacheKey($iterationId);

        if ($cache = Cache::store('file')->get($iterationListCacheKey)) {
            $cache[$iterationCacheKey] = $tree;
        } else {
            $cache = [
                $iterationCacheKey => $tree
            ];
        }

        return Cache::store('file')->put($iterationListCacheKey, $cache, now()->addDay());
    }

    /**
     * 删除项目相关所有文档树的缓存
     *
     * @param int $projectID 项目id
     * @return void
     */
    public static function remove(int $projectID)
    {
        $projectCacheKey = self::projectCacheKey($projectID);
        $iterationListCacheKey = self::iterationListCacheKey($projectID);

        Cache::store('file')->forget($projectCacheKey);
        Cache::store('file')->forget($iterationListCacheKey);
    }

    /**
     * 获取项目缓存key
     *
     * @param int $projectId 项目id
     * @return string
     */
    public static function projectCacheKey(int $projectId)
    {
        return 'project_tree:' . $projectId;
    }

    /**
     * 获取项目迭代缓存key
     *
     * @param int $projectId 项目id
     * @return string
     */
    public static function iterationListCacheKey(int $projectId)
    {
        return 'iteration_tree:' . $projectId;
    }

    /**
     * 获取迭代缓存key
     *
     * @param int $iterationId 迭代id
     * @return string
     */
    public static function iterationCacheKey(int $iterationId)
    {
        return 'tree:' . $iterationId;
    }
}
