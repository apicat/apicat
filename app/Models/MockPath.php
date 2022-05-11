<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;

class MockPath extends Model
{
    use HasFactory, SoftDeletes;

    protected $fillable = [
        'project_id',
        'doc_id',
        'path',
        'format',
        'method'
    ];
}
