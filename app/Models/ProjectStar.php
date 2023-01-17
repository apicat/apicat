<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class ProjectStar extends Model
{
    use HasFactory;

    protected $fillable = [
        'user_id', 'project_id', 'display_order',
    ];
}
