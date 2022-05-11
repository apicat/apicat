<?php

namespace App\Repositories\Project;

use App\Models\Project;
use App\Models\ProjectMember;
use Illuminate\Support\Facades\Auth;

class ActiveProjectRepository
{
    /**
     * 项目实例
     *
     * @var Project
     */
    public $project;

    /**
     * 成员在项目中的权限
     *
     * @var int|null
     */
    public $authority;

    public function __construct(Project $project)
    {
        $this->project = $project;
    }

    public function __get($name)
    {
        if (isset($this->project->$name)) {
            return $this->project->$name;
        }
        return null;
    }

    /**
     * 检查成员是否有文档编辑权限
     *
     * @return boolean
     */
    public function hasAuthority()
    {
        if (is_null($this->authority)) {
            $this->authority = ProjectMember::where([
                ['project_id', $this->project->id],
                ['user_id', Auth::id()]
            ])->value('authority');
        }

        return 2 > $this->authority;
    }

    /**
     * 检查成员是否是项目管理者
     * 
     * @return boolean 是: true  不是: false
     */
    public function isAdmin()
    {
        if (is_null($this->authority)) {
            $this->authority = ProjectMember::where([
                ['project_id', $this->project->id],
                ['user_id', Auth::id()]
            ])->value('authority');
        }

        return 0 === $this->authority;
    }
}
